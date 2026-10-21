//go:build wolfssl

package profile

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/florianoneto07/akma-free5gc/af/af/internal/config"
	"github.com/florianoneto07/akma-free5gc/af/af/internal/logger"
	"github.com/florianoneto07/akma-free5gc/af/af/internal/server"
	"github.com/florianoneto07/akma-free5gc/af/akmatls"
	"github.com/florianoneto07/akma-free5gc/af/wolftls"
	"github.com/florianoneto07/akma-free5gc/shared/akma/afid"
)

// Built is true: this AF has the PSK profiles, and with wolfSSL it is
// licensed under the GNU General Public License v3.0.
const Built = true

// pskListener returns the listener of a PSK profile (TS 24.109 K.3.2 and
// K.3.3), the ConnContext that gives the handler the session of its
// connection, and the release of the wolfSSL context.
func pskListener(p config.Profile, s *config.Settings, ua *server.Server, l net.Listener) (net.Listener,
	func(context.Context, net.Conn) context.Context, func(), error) {
	version, name := wolftls.VersionTLS13, "TLS 1.3"
	if p.Type == akmatls.ProfilePSKTLS12 {
		version, name = wolftls.VersionTLS12, "TLS 1.2"
	}
	ws, err := wolftls.NewServer(wolftls.ServerConfig{
		Version:         version,
		ServerName:      s.FQDN,
		IdentityHint:    akmatls.IdentityHint,
		SessionLifetime: akmatls.SessionLifetime,
		Key: func(req wolftls.KeyRequest) (wolftls.Key, error) {
			// TLS 1.3: the fixed identifier, known to the UE before the
			// ClientHello; TLS 1.2: the one of the selected cipher suite
			// (TS 33.222 cl. 5.4.0.2 NOTE 0; TS 33.220 H.3).
			id := afid.AfID{FQDN: s.FQDN, Protocol: afid.UaStarTLS13}
			if version == wolftls.VersionTLS12 {
				id.Protocol = akmatls.UaStarTLS(req.CipherSuite)
			}
			ctx, cancel := context.WithTimeout(context.Background(), config.DefaultAanfTimeout+time.Second)
			defer cancel()
			sess, kaf, err := ua.PSKKey(ctx, id, req.Identity, req.RemoteAddr.String())
			if err != nil {
				return wolftls.Key{}, err
			}
			// Expiry on the clock of the process, not on a shifted one.
			return wolftls.Key{PSK: kaf[:], Expiry: sess.KafExpiry.Add(-s.ClockOffset), Peer: sess}, nil
		},
		Expired: func(c *wolftls.Conn) {
			if sess, ok := c.State().Peer.(server.PSKSession); ok {
				logger.UaLog.Infof("K_AF of %s expired: TLS session with %s closed", sess, c.RemoteAddr())
			}
		},
	})
	if err != nil {
		return nil, nil, nil, err
	}
	connCtx := func(ctx context.Context, c net.Conn) context.Context {
		lc, ok := c.(*loggedConn)
		if !ok {
			return ctx
		}
		return server.WithPSKSession(ctx, func() (server.PSKSession, bool) {
			sess, ok := lc.State().Peer.(server.PSKSession)
			return sess, ok
		})
	}
	return &loggingListener{Listener: ws.Listener(l), name: name}, connCtx, ws.Close, nil
}

// loggingListener logs the outcome of each handshake, which net/http does
// not see on a connection that is not a *tls.Conn.
type loggingListener struct {
	net.Listener
	name string
}

func (l *loggingListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &loggedConn{Conn: c.(*wolftls.Conn), name: l.name}, nil
}

type loggedConn struct {
	*wolftls.Conn
	name string
	once sync.Once
}

func (c *loggedConn) Read(p []byte) (int, error) {
	c.once.Do(func() {
		if err := c.Handshake(); err != nil {
			logger.UaLog.Warnf("%s PSK handshake with %s failed: %v", c.name, c.RemoteAddr(), err)
			return
		}
		st := c.State()
		sess, _ := st.Peer.(server.PSKSession)
		logger.UaLog.Infof("%s PSK handshake with %s: %s, cipher suite %04X, resumed %t", c.name, c.RemoteAddr(),
			sess, st.CipherSuite, st.Resumed)
	})
	return c.Conn.Read(p)
}

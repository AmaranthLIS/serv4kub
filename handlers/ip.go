package handlers

import (
    "net/http"
    "net"
    "github.com/rs/zerolog/log"
)

func GetClientIP() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(ipProxy))
            return
        }
        ip, _, _ := net.SplitHostPort(r.RemoteAddr)
        log.Info().Str("ask", ip).Msg("")

        w.WriteHeader(http.StatusOK)
        w.Write([]byte(ip))
    }
}

func GetServerIp(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(GetLocalIP()))
}

func GetLocalIP() string {
    addr, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addr {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

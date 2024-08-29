package ipresolver

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Тут мы смотрим на основе чего доставать IP
type ResolveIPOpts struct {
	UseHeader bool
}

// ResolveIP - получаем IP из заголовка или иного
func ResolveIP(r *http.Request, opts ResolveIPOpts) (net.IP, error) {
	if !opts.UseHeader {
		addr := r.RemoteAddr
		// метод возвращает адрес в формате host:port
		// нужна только подстрока host
		ipStr, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		// парсим ip
		ip := net.ParseIP(ipStr)
		if ip == nil {
			panic("unexpected parse ip error")
		}
		return ip, nil
	} else {
		// смотрим заголовок запроса X-Real-IP
		ipStr := r.Header.Get("X-Real-IP")
		// парсим ip
		ip := net.ParseIP(ipStr)
		if ip == nil {
			// если заголовок X-Real-IP пуст, пробуем X-Forwarded-For
			// этот заголовок содержит адреса отправителя и промежуточных прокси
			// в виде 203.0.113.195, 70.41.3.18, 150.172.238.178
			ips := r.Header.Get("X-Forwarded-For")
			// разделяем цепочку адресов
			ipStrs := strings.Split(ips, ",")
			// интересует только первый
			ipStr = ipStrs[0]
			// парсим
			ip = net.ParseIP(ipStr)
		}
		if ip == nil {
			return nil, fmt.Errorf("failed parse ip from http header")
		}
		return ip, nil
	}
}

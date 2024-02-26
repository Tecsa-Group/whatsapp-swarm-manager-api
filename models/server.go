package models

import "time"

type Server struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	URL       string    `json:"url"`
}

type RequestServerHetzner struct {
	Firewalls        []Firewall `json:"firewalls"`
	Image            string     `json:"image"`
	Name             string     `json:"name"`
	ServerType       string     `json:"server_type"`
	SSHKeys          []int      `json:"ssh_keys"`
	StartAfterCreate bool       `json:"start_after_create"`
}

type Firewall struct {
	Firewall int `json:"firewall"`
}

type ResponseServerHetzner struct {
	Server ServerHetzner `json:"server"`
}

type ServerHetzner struct {
	PublicNet PublicNet `json:"public_net"`
	Name      string    `json:"name"`
}

type PublicNet struct {
	IPv4 IPv4 `json:"ipv4"`
}

type IPv4 struct {
	IP      string `json:"ip"`
	DNSPtr  string `json:"dns_ptr"`
	Blocked bool   `json:"blocked"`
}

type DNSRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl,omitempty"`
	Proxied bool   `json:"proxied,omitempty"`
}

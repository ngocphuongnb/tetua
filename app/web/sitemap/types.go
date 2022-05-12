package websitemap

import "encoding/xml"

const SITEMAP_PAGESIZE = 10000

type SitemapItems struct {
	XMLName xml.Name       `xml:"sitemapindex"`
	Xmlns   string         `xml:"xmlns,attr"`
	Items   []*SitemapItem `xml:"sitemap"`
}

type SitemapItem struct {
	Loc     *SitemapLoc
	LastMod *SitemapLastMod
}

type SitemapLoc struct {
	XMLName xml.Name `xml:"loc"`
	Value   string   `xml:",chardata"`
}

type SitemapLastMod struct {
	XMLName xml.Name `xml:"lastmod"`
	Value   string   `xml:",chardata"`
}

type SitemapUrlSets struct {
	XMLName             xml.Name      `xml:"urlset"`
	Xmlns               string        `xml:"xmlns,attr"`
	XmlnsXsi            string        `xml:"xmlns:xsi,attr"`
	XmlnsImage          string        `xml:"xmlns:image,attr"`
	XmlnsSchemaLocation string        `xml:"xsi:schemaLocation,attr"`
	Urls                []*SitemapUrl `xml:"url"`
}

type SitemapUrl struct {
	Loc     *SitemapLoc
	LastMod *SitemapLastMod
	Image   *SitemapImage
}

type SitemapImage struct {
	XMLName xml.Name `xml:"image:image"`
	Loc     *SitemapImageLoc
}

type SitemapImageLoc struct {
	XMLName xml.Name `xml:"image:loc"`
	Value   string   `xml:",chardata"`
}

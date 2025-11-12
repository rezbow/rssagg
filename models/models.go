package models

type FeedItem struct {
	ID      int
	Title   string
	Content string
}

type Feed struct {
	ID    int
	Title string
	Items []FeedItem
}

/*****************************************************************************
 **
 ** grr >:(
 ** https://github.com/melllvar/grr
 ** Copyright (C) 2013 Akop Karapetyan
 **
 ** This program is free software; you can redistribute it and/or modify
 ** it under the terms of the GNU General Public License as published by
 ** the Free Software Foundation; either version 2 of the License, or
 ** (at your option) any later version.
 **
 ** This program is distributed in the hope that it will be useful,
 ** but WITHOUT ANY WARRANTY; without even the implied warranty of
 ** MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 ** GNU General Public License for more details.
 **
 ** You should have received a copy of the GNU General Public License
 ** along with this program; if not, write to the Free Software
 ** Foundation, Inc., 675 Mass Ave, Cambridge, MA 02139, USA.
 **
 ******************************************************************************
 */
 
package parser

import (
  "time"
  "encoding/xml"
)

var supportedAtomTimeFormats = []string {
  "2006-01-02T15:04:05Z07:00",
}

type atomFeed struct {
  XMLName xml.Name `xml:"feed"`
  Id string `xml:"id"`
  Title string `xml:"title"`
  Description string `xml:"subtitle"`
  Updated string `xml:"updated"`
  Link []atomLink `xml:"link"`
  Entry []*atomEntry `xml:"entry"`
}

type atomLink struct {
  Type string `xml:"type,attr"`
  Rel string `xml:"rel,attr"`
  Href string `xml:"href,attr"`
}

type atomAuthor struct {
  Name string `xml:"name"`
  URI string `xml:"uri"`
}

type atomEntry struct {
  Id string `xml:"id"`
  Published string `xml:"published"`
  Updated string `xml:"updated"`
  Link []atomLink `xml:"link"`
  EntryTitle atomText `xml:"title"`
  Content atomText `xml:"content"`
  Summary atomText `xml:"summary"`
  Author atomAuthor `xml:"author"`
}

type atomText struct {
  Type string `xml:"type,attr"`
  Content string `xml:",chardata"`
}

func (nativeFeed *atomFeed) Marshal() (feed Feed, err error) {
  updated := time.Time {}
  if nativeFeed.Updated != "" {
    updated, err = parseTime(supportedAtomTimeFormats, nativeFeed.Updated)
  }

  linkUrl := ""
  for _, link := range nativeFeed.Link {
    if link.Rel == "alternate" {
      linkUrl = link.Href
    }
  }

  feed = Feed {
    Title: nativeFeed.Title,
    Description: nativeFeed.Description,
    Updated: updated,
    WWWURL: linkUrl,
  }

  if nativeFeed.Entry != nil {
    feed.Entry = make([]*Entry, len(nativeFeed.Entry))
    for i, v := range nativeFeed.Entry {
      var entryError error
      feed.Entry[i], entryError = v.Marshal()

      if entryError != nil && err == nil {
        err = entryError
      }
    }
  }

  return feed, err
}

func (nativeEntry *atomEntry) Marshal() (entry *Entry, err error) {
  linkUrl := ""
  for _, link := range nativeEntry.Link {
    if link.Rel == "alternate" {
      linkUrl = link.Href
    }
  }

  guid := nativeEntry.Id
  if guid == "" {
    guid = linkUrl
  }

  content := nativeEntry.Content.Content
  if content == "" && nativeEntry.Summary.Content != "" {
    content = nativeEntry.Summary.Content
  }

  published := time.Time {}
  if nativeEntry.Published != "" {
    published, err = parseTime(supportedAtomTimeFormats, nativeEntry.Published)
  }

  updated := time.Time {}
  if nativeEntry.Updated != "" {
    updated, err = parseTime(supportedAtomTimeFormats, nativeEntry.Updated)
  }

  entry = &Entry {
    GUID: guid,
    Author: nativeEntry.Author.Name,
    Title: nativeEntry.EntryTitle.Content,
    Content: content,
    Published: published,
    Updated: updated,
    WWWURL: linkUrl,
  }

  return entry, err
}

# booklog-tool
A CLI tool to update book information on booklog.jp

## Setup
Create config file `~/.config/booklog-tool/config.json`.  You can obtain the Cookie (`bid`) from browsers' developer tools after you sign in to booklog. 
```json
{
    "username": "your_username",
    "cookie": "your_bid"
}
```

## Install
```sh
# Install the binary globally
go install github.com/mu373/booklog-tool@HEAD
booklog-tool

# Or clone the repository and build it yourself locally
go build .
./booklog-tool
```

## Usage
Prepare a text file containing booklog item IDs.
```txt
4051331255
4051331256
4051331257
```

```sh
# Add a new tag to item
./booklog-tool add-tag -t "your_new_tag" -f list.txt

# Update location tag (e.g., "loc_Tokyo", "loc_Osaka", "loc_London") of the item
# If a location tag already exists, it will be overrided with the new location tag
./booklog-tool update-location -l "Osaka" -f list.txt
```

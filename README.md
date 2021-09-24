# CocDiscordLinkGo
Implementation of reverendmike's CocDiscordLinkAPI for Go

# Usage
    import "github.com/eslindsey/CocDiscordLinkGo"
    
    const (
        Username = "The username assigned to your project by reverendmike"
        Password = "The password assigned to your project by reverendmike"
    )
    
    func main() {
        s, err := cocdiscordlink.New(username, password)
    }

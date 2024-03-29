package bot

import (
	"strings"
	"time"

	"github.com/aiuzu42/SukiBot/config"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func CommandsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.GuildID == "" {
		handleDM(s, m)
		return
	}

	if strings.HasPrefix(m.Content, prefix) == true {
		r := []rune(m.Content)
		st := string(r[pLen:])
		args := strings.Split(st, " ")
		for _, custom := range config.Config.CustomSays {
			if custom.Name == args[0] {
				customSayCommand(s, m, custom.Channel, st)
				return
			}
		}
		switch args[0] {
		case "reloadConfig":
			reloadConfigCommand(s, m.ChannelID, m.ID, m.Author.ID)
		case "setStatus":
			setGameStatus(s, m, args)
		case "setListenStatus":
			setListenStatus(s, m, args)
		case "setStreamStatus":
			setStreamStatus(s, m, args)
		}
	} else {
		for _, t := range config.Config.Triggers {
			if !t.AllowsChannel(m.ChannelID) {
				continue
			}
			if (t.CaseSensitive && strings.Contains(m.Content, t.Trigger)) || (!t.CaseSensitive && strings.Contains(strings.ToLower(m.Content), t.Trigger)) {
				triggerResponse(s, m.ChannelID, t)
				break
			}
		}
	}
}

func customSayCommand(s *discordgo.Session, m *discordgo.MessageCreate, tarCh string, st string) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[customSayCommand]User: " + m.Author.ID + " tried to use command customSayCommand without permission.")
		return
	}
	sayCommand(s, m.ChannelID, tarCh, m.ID, st)
}

func sayCommand(s *discordgo.Session, originCh string, tarCh string, id string, st string) {
	err := s.ChannelMessageDelete(originCh, id)
	if err != nil {
		log.Error("[sayCommand]Can't delete message: " + err.Error())
	}
	msg := saySplit(st)
	_, err = s.ChannelMessageSend(tarCh, msg)
	if err != nil {
		log.Error("[sayCommand]Can't send message: " + err.Error())
	}
}

func reloadConfigCommand(s *discordgo.Session, channelID string, messageId string, id string) {
	if !IsOwner(id) {
		log.Warn("[reloadRolesCommand]User: " + id + " tried to use command reloadRolesCommand without permission.")
		return
	}
	if err := config.ReloadConfig(); err != nil {
		log.Error("[reloadRolesCommand]Error reloading config: " + err.Error())
		sendMessage(s, channelID, "Hubo un error, no se pudo recargar la congfiguracion", "[reloadRolesCommand][1]")
	}
	LoadRoles()
	err := s.MessageReactionAdd(channelID, messageId, "✅")
	if err != nil {
		log.Error("[reloadRolesCommand]Error marking message: " + err.Error())
	}
}

func handleDM(s *discordgo.Session, m *discordgo.MessageCreate) {
}

func saySplit(st string) string {
	msg := ""
	args := strings.Split(st, " ")
	msg = strings.Join(args[1:], " ")
	return msg
}

func sendMessage(s *discordgo.Session, channelID string, msg string, logMsg string) {
	_, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		log.Error(logMsg + "Error sending message [" + msg + "]: " + err.Error())
	}
}

func sendSimpleEmbedMessage(s *discordgo.Session, channelID string, msg string, color int, image string, logMsg string) {
	me := discordgo.MessageEmbed{Description: msg, Color: color}
	if image != "" {
		me.Image = &discordgo.MessageEmbedImage{URL: image}
	}
	_, err := s.ChannelMessageSendEmbed(channelID, &me)
	if err != nil {
		log.Error(logMsg + "Error sending message [" + msg + "]: " + err.Error())
	}
}

func setGameStatus(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[setGameStatus]User: " + m.Author.ID + " tried to use command setGameStatus without permission.")
		return
	}
	if len(args) < 2 {
		log.Error("[setGameStatus][1]Invalid number of arguments")
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, favor de revisar el comando", "[setGameStatus][1]")
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("[setGameStatus]Unable to delete message: " + err.Error())
	}
	msg := strings.Join(args[1:], " ")
	err = s.UpdateGameStatus(0, msg)
	if err != nil {
		log.Error("[setGameStatus]Unable to update status: " + err.Error())
	}
}

func setListenStatus(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[setListenStatus]User: " + m.Author.ID + " tried to use command setListenStatus without permission.")
		return
	}
	if len(args) < 2 {
		log.Error("[setListenStatus][1]Invalid number of arguments")
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, favor de revisar el comando", "[setListenStatus][1]")
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("[setListenStatus]Unable to delete message: " + err.Error())
	}
	msg := strings.Join(args[1:], " ")
	err = s.UpdateListeningStatus(msg)
	if err != nil {
		log.Error("[setListenStatus]Unable to update status: " + err.Error())
	}
}

func setStreamStatus(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[setStreamStatus]User: " + m.Author.ID + " tried to use command setStreamStatus without permission.")
		return
	}
	if len(args) < 3 {
		log.Error("[setStreamStatus][1]Invalid number of arguments")
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, favor de revisar el comando", "[setStreamStatus][1]")
		return
	}
	url := args[1]
	msg := strings.Join(args[2:], " ")
	if url == "" && msg == "" {
		log.Error("[setStreamStatus][2]Invalid number of arguments")
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, favor de revisar el comando", "[setStreamStatus][2]")
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("[setStreamStatus]Unable to delete message: " + err.Error())
	}
	err = s.UpdateStreamingStatus(0, msg, url)
	if err != nil {
		log.Error("[setStreamStatus]Unable to update status: " + err.Error())
	}
}

func triggerResponse(s *discordgo.Session, ch string, t config.Trigger) {
	now := time.Now()
	r, ok := t.CooldownMap[ch]
	if !ok || (ok && r+t.Cooldown <= now.Unix()) {
		sendSimpleEmbedMessage(s, ch, t.Response, t.Color, t.Image, "[triggerResponse]")
		t.CooldownMap[ch] = now.Unix()
	} else {
		log.Infof("[triggerResponse]Trigger in timer, now: %d, map: %d, cd: %d", now.Unix(), r, t.Cooldown)
	}
}

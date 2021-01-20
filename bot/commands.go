package bot

import (
	"strings"

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
			setStatus(s, m, r)
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

func setStatus(s *discordgo.Session, m *discordgo.MessageCreate, r []rune) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[setStatus]User: " + m.Author.ID + " tried to use command setStatus without permission.")
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("[setStatus]Unable to delete message: " + err.Error())
	}
	msg := string(r[pLen+10:])
	err = s.UpdateStatus(0, msg)
	if err != nil {
		log.Error("[setStatus]Unable to update status: " + err.Error())
	}
}

import gql from 'graphql-tag'

export const REMOVE_FILE = gql`mutation ($file:String!) {
  removeFile(file: $file)
}`

export const SET_MAPPING = gql`mutation ($slot: Int!, $name:String!) {
  setSlotMapping(slot: $slot, name:$name)
}`

export const REMOVE_SLOT = gql`mutation ($slot:Int!) {
  removeSlotMapping(slot: $slot)
}`

export const SET_PLAYING = gql`mutation ($slot:Int!, $player:PlayerInput!) {
  setPlaying(slot:$slot, player:$player) {
    playing
    loop
    volume
  }
}`

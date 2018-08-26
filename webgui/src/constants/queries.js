import gql from 'graphql-tag'

export const GET_SLOTS = gql`{
  slots{
    slot
    item {
      name
      extension
    }
    player {
      playing
      loop
      volume
      rate
      current
      length
    }
  }
}`

export const GET_POOL = gql`{
  pool{
    name
    extension
    used
  }
}`

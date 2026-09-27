import { ExchangeListItems } from "./ExchangeListItem"

// TODO: all types need to be shared from the db types

export type Exchange = {
    id: string,
    name: string,
    exchangeDate: Date,
    exchangeDescription?: string,
    exchangeBudgent?: number,
    exchangeInites: string[],
    exchangeOrganizerId: string
}

type Props = {
    exchanges: Exchange[]
}


export function ExchangeList({ exchanges }: Props) {
    return <ul>
        {exchanges.map(e => <ExchangeListItems key={e.id} {...e} />)}
    </ul>
}
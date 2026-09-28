import type { Exchange } from "@giftexchanger/types"
import { ExchangeListItems } from "./ExchangeListItem"

type Props = {
    exchanges: Exchange[]
}


export function ExchangeList({ exchanges }: Props) {
    return <ul>
        {exchanges.map(e => <ExchangeListItems key={e.id} {...e} />)}
    </ul>
}

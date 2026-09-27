'use client'
import { Button } from '@/components/ui/button'
import { Gift } from 'lucide-react'
import { useSearchParams } from 'next/navigation'
import { use } from 'react'

// /invate/abc123?email=matthew.william.ritter@gmail.com

export default function InvitePage({
    params,
}: {
    params: Promise<{ exchangeId: string }>
}) {
    const { exchangeId } = use(params)
    const query = useSearchParams()
    const email = query.get('email')

    // Fetch the exchange to get the name by this id

    const onJoinExchange = () => {
        // the backend will need to check if the email is in the exchange
        // then accept exchange invite
        // I wonder if the invite link should be a magic link so if they click the email thats creates the 'token' account
        // and when they click join here they accept the invite and are redirected to /dashboard/exchanges
        console.log({
            exchangeId,
            email
        })
    }

    return (
        <div className='relative flex flex-col justify-center items-center h-screen bg-muted'>
            <div className='flex flex-col gap-4 md:px-20 md:pb-5 p-10 bg-white rounded-md shadow-xl'>
                <div className='flex flex-col items-center -mt-4'>
                    <Gift size={50} className='rounded-full p-2 bg-muted' />
                    <p className='font-bold'>GiftExchange</p>
                </div>
                <div className='flex-1 flex flex-col items-center gap-5'>
                    <h1 className='text-center text-2xl font-bold'>Join a Gift Exchange</h1>
                    <div>
                        <p>You&apos;ve been invited to join Exchange Name!</p>
                        <p>Click Join exchange below to get started.</p>
                    </div>
                </div>
                <Button className="mt-10 text-md" size="lg" onClick={onJoinExchange}>Join exchange</Button>
            </div>
        </div>
    )
}
import type {Secret} from "../../util/secret.ts";

type SecretListProps = {
    secrets: Secret[]
    clickedSecret: (s: Secret) => void
};

export const SecretsList = (props: SecretListProps) => {
    return <ul className="list bg-base-100 rounded-box shadow-md w-full">

        {
            props.secrets.length == 0
                ? <li className="p-4 pb-2 text-l opacity-60 tracking-wide">No secrets..</li>
                : <li className="p-4 pb-2 text-l opacity-60 tracking-wide">Your secrets:</li>
        }

        {props.secrets.map((secret) => (
            <li className="list-row w-full" key={secret.id}>
                <button className="btn btn-ghost w-full flex flex-col min-h-fit lg:text-lg p-1"
                        onClick={() => props.clickedSecret(secret)}>
                    <div className="flex flex-col max-w-full items-start">
                        <p className="text-semibold truncate max-w-full">{secret.key}</p>
                        <p className="font-normal italic truncate max-w-full">[{secret.url}]</p>
                        <div className="flex flex-row gap-2 max-w-full items-center truncate pt-1.5">
                            {secret.tags.map((tag) => (
                                <span className="badge badge-neutral badge-xs">{tag}</span>
                            ))}
                        </div>
                    </div>
                </button>
            </li>
        ))}
    </ul>
}

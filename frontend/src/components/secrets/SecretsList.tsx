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
                <button className="btn btn-ghost min-w-full justify-start" onClick={() => props.clickedSecret(secret)}>
                    <div className="text-xl flex flex-row">
                        <p className="font-semibold">{secret.key}</p>
                        &nbsp;
                        <p className="font-normal italic">[{secret.url}]</p>
                    </div>
                </button>
            </li>
        ))}
    </ul>
}

import type {UnlockRequest} from "../../util/api.ts";

export const UnlockRequests = ({requests, contribute}: {requests: UnlockRequest[], contribute: (id: number) => Promise<void>}) => {
    if (requests.length === 0) return null;
    return <div className="flex flex-col gap-2">
        {requests.map(request => <div role="status" className="alert alert-warning" key={request.id}>
            <div className="grow">
                <div className="font-semibold">{request.requesterName} is unlocking “{request.secretName}”</div>
                <div className="text-sm">{request.contributions} of {request.threshold} shares contributed</div>
            </div>
            <button className="btn btn-sm btn-primary" disabled={request.contributed} onClick={() => contribute(request.id)}>
                {request.contributed ? "Share contributed" : "Participate"}
            </button>
        </div>)}
    </div>;
};

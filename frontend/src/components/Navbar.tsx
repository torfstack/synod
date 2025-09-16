import ThemeSwitcher from "./ThemeSwitcher.tsx";
import {useEffect, useState} from "react";
import {deleteAuth, getAuth, getSecrets, postSecret} from "../util/api.ts";
import {config} from "../util/config.ts";
import Hero from "./Hero.tsx";
import SecretList from "./secrets/SecretList.tsx";
import {SecretModal, showModal} from "./secrets/SecretModal.tsx";
import type {Secret} from "../util/secret.ts";

const Navbar: React.FC = () => {
  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    getAuth().then(res => {
      if (res.status == 200) {
        setAuthenticated(true);
      }
    })
  }, [])

  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [selectedSecret, setSelectedSecret] = useState<Secret | undefined>(undefined);

  async function uploadSecret(s: Secret) {
    await postSecret(s)
    retrieveSecrets()
  }

  function selectSecret(s: Secret) {
    setSelectedSecret(s)
    showModal()
  }

  function retrieveSecrets() {
    getSecrets().then(value => value.json()).then(
      (json) => {
        setSecrets(json)
      }
    )
  }

  useEffect(() => {
    retrieveSecrets()
  }, [])

  return (
    <div className="flex flex-col h-screen">
      <div className="navbar bg-base-100 shadow-sm">
        <div className="flex-1">
          <a className="btn btn-ghost text-l" href={"/"}>
            <span className="badge">
              KayVault
            </span>
          </a>
        </div>


        <div className="flex gap-4 items-center">
          {authenticated && <button className={"btn"} onClick={() => {
            const emptySecret: Secret = {
              id: undefined, key: undefined, tags: [], url: undefined, value: undefined
            }
            setSelectedSecret(emptySecret);
            showModal()
          }}>Add Secret</button>}

          {
            authenticated &&
            <input
              type="text"
              placeholder="Search"
              className="input input-bordered w-36 md:w-auto"
            />
          }

          {
            authenticated
              ? <LogoutButton setAuthenticated={setAuthenticated}/>
              : <LoginButton/>
          }

          <ThemeSwitcher/>
        </div>
      </div>

      {
        authenticated
          ? <div className="h-full bg-base-200">
            <div className="justify-center p-8">
              <SecretList secrets={secrets} clickedSecret={selectSecret}/>
            </div>
            <SecretModal handleSecret={uploadSecret} selectedSecret={selectedSecret}/>
          </div>
          : <div className="h-full bg-base-200 flex flex-col justify-center"><Hero/></div>
      }
    </div>
  );
};

type SetAuthenticatedProps = {
  setAuthenticated: (value: boolean) => void
}

const LoginButton = () => {
  const signInWithProvider = () => {
    window.open(config.backendAuthStartUrl, "_self");
  }

  return <button className="btn btn-primary" onClick={signInWithProvider}>
    <svg aria-label="Google logo" width="16" height="16" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
      <g>
        <path fill="#34a853" d="M153 292c30 82 118 95 171 60h62v48A192 192 0 0190 341"></path>
        <path fill="#4285f4" d="m386 400a140 175 0 0053-179H260v74h102q-7 37-38 57"></path>
        <path fill="#fbbc02" d="m90 341a208 200 0 010-171l63 49q-12 37 0 73"></path>
        <path fill="#ea4335" d="m153 219c22-69 116-109 179-50l55-54c-78-75-230-72-297 55"></path>
      </g>
    </svg>
    Login with Google
  </button>
}

const LogoutButton = (props: SetAuthenticatedProps) => {
  async function signOutWithProvider() {
    await deleteAuth()
    props.setAuthenticated(false);
  }

  return <button className="btn btn-accent" onClick={signOutWithProvider}>
    Logout
  </button>
}

export default Navbar;
export function showPasswordPrompt(): Promise<string> {
  return new Promise<string>((resolve, reject) => {
    const promptBackdrop = document.createElement('div');
    promptBackdrop.classList.add('prompt-backdrop');
    const passwordPrompt = document.createElement('div');
    passwordPrompt.classList.add('password-prompt');
    passwordPrompt.innerHTML=`
        <div class="row">
            <div class="col-lg-12">
                <div class="form-group">
                    <label>Provide password</label>
                    <input class="form-control" type="password" />
                </div>
            </div>
        </div>
        <div class="row mt-2">
            <div class="col-lg-12">
                <button class="btn btn-primary">Ok</button>
            </div>
        </div>
      `;
    const login = () => {
      document.querySelector('body').removeChild(passwordPrompt);
      document.querySelector('body').removeChild(promptBackdrop);
      resolve(passwordPrompt.querySelector('input').value);
    };
    passwordPrompt.querySelector('button').addEventListener('click', login);
    passwordPrompt.querySelector('input').addEventListener('keydown', (event) => {
      if(event.key === 'Enter') {
        login();
      }
    });
    document.querySelector('body').appendChild(passwordPrompt);
    document.querySelector('body').appendChild(promptBackdrop);
  });
}

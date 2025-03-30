
document.addEventListener('DOMContentLoaded', function () {
    const steps = document.querySelectorAll('.form-step');
    const nextButtons = document.querySelectorAll('.next-step');
    const prevButtons = document.querySelectorAll('.prev-step');
    const progressBar = document.querySelector('.progress-bar');
    const stepIndicators = document.querySelectorAll('.step');
    let currentStep = 0;

    // Helper function for showing errors
    async function showError(title, message) {
        await Swal.fire({
            icon: 'error',
            title: title,
            text: message,
            confirmButtonColor: '#007bff'
        });
    }

    // Helper function for showing success
    async function showSuccess(title, message) {
        await Swal.fire({
            icon: 'success',
            title: title,
            text: message,
            confirmButtonColor: '#007bff'
        });
    }

    // Helper function for showing warnings
    async function showWarning(title, message) {
        await Swal.fire({
            icon: 'warning',
            title: title,
            text: message,
            confirmButtonColor: '#007bff'
        });
    }

    // Xử lý nút "Next"
    nextButtons.forEach(button => {
        button.addEventListener('click', async function (e) {
            e.preventDefault();

            let isValid = false;
            switch (currentStep) {
                case 0:
                    isValid = await validateStep1();
                    break;
                case 1:
                    isValid = await validateStep2();
                    break;
            }

            if (!isValid) {
                return;
            }

            steps[currentStep].classList.remove('active');
            currentStep++;
            steps[currentStep].classList.add('active');

            progressBar.style.width = `${(currentStep + 1) * 33.33}%`;

            stepIndicators[currentStep - 1].classList.add('completed');
            stepIndicators[currentStep].classList.add('active');
        });
    });

    prevButtons.forEach(button => {
        button.addEventListener('click', function () {
            steps[currentStep].classList.remove('active');
            stepIndicators[currentStep].classList.remove('active');
            currentStep--;
            steps[currentStep].classList.add('active');

            progressBar.style.width = `${(currentStep + 1) * 33.33}%`;
        });
    });

    function handleImagePreview(inputId, previewId, placeholderId) {
        document.getElementById(inputId).addEventListener('change', function (e) {
            const file = e.target.files[0];
            if (file) {
                const reader = new FileReader();
                reader.onload = function (e) {
                    const preview = document.getElementById(previewId);
                    preview.src = e.target.result;
                    preview.style.display = 'block';
                    document.getElementById(placeholderId).style.display = 'none';
                }
                reader.readAsDataURL(file);
            }
        });
    }

    handleImagePreview('idCardFront', 'frontIdPreview', 'frontIdPlaceholder');
    handleImagePreview('idCardBack', 'backIdPreview', 'backIdPlaceholder');

    // Gửi form đăng ký
    document.getElementById('registerForm').addEventListener('submit', async function (event) {
        event.preventDefault();

        if (currentStep !== steps.length - 1) {
            return;
        }

        if (!await validateStep3()) {
            return;
        }

        const submitButton = this.querySelector('button[type="submit"]');
        submitButton.disabled = true;

        try {
            const formData = {
                username: document.getElementById('identifier').value,
                email: document.getElementById('email').value,
                password: document.getElementById('password').value,
                firstName: document.getElementById('firstName').value,
                lastName: document.getElementById('lastName').value,
                phoneNumber: document.getElementById('phone').value,
                address: document.getElementById('address').value,
                city: document.getElementById('city').value,
                zipCode: document.getElementById('zipCode').value,
            };

            // Handle ID card images
            const frontIdFile = document.getElementById('idCardFront').files[0];
            const backIdFile = document.getElementById('idCardBack').files[0];

            if (frontIdFile) {
                formData.idCardFront = await fileToBase64(frontIdFile);
            }
            if (backIdFile) {
                formData.idCardBack = await fileToBase64(backIdFile);
            }

            console.log('Sending registration data...');

            // Show loading state
            Swal.fire({
                title: 'Processing...',
                text: 'Please wait while we create your account',
                allowOutsideClick: false,
                showConfirmButton: false,
                willOpen: () => {
                    Swal.showLoading();
                }
            });

            const response = await fetch('/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Registration failed');
            }

            const data = await response.json();
            console.log('Registration successful:', data);

            // Show success message
            await showSuccess('Registration Successful', 'Your account has been created successfully! Redirecting to login page...');

            // Redirect after success message
            window.location.href = '/login';

        } catch (error) {
            console.error('Registration error:', error);
            await showError('Registration Failed', error.message);
        } finally {
            submitButton.disabled = false;
        }
    });

    async function validateStep1() {
        const username = document.getElementById('identifier').value;
        const password = document.getElementById('password').value;
        const confirmPassword = document.getElementById('confirmPassword').value;
        const firstName = document.getElementById('firstName').value;
        const lastName = document.getElementById('lastName').value;

        // Username validation
        if (username.length < 3 || username.length > 30) {
            await showWarning('Invalid Username', 'Username must be between 3 and 30 characters');
            return false;
        }
        if (!/^[a-zA-Z0-9_]+$/.test(username)) {
            await showWarning('Invalid Username', 'Username can only contain letters, numbers, and underscores');
            return false;
        }

        // Password validation
        if (password.length < 8) {
            await showWarning('Weak Password', 'Password must be at least 8 characters long');
            return false;
        }
        if (!/[A-Z]/.test(password)) {
            await showWarning('Weak Password', 'Password must contain at least one uppercase letter');
            return false;
        }
        if (!/[a-z]/.test(password)) {
            await showWarning('Weak Password', 'Password must contain at least one lowercase letter');
            return false;
        }
        if (!/[0-9]/.test(password)) {
            await showWarning('Weak Password', 'Password must contain at least one number');
            return false;
        }
        if (!/[!@#$%^&*]/.test(password)) {
            await showWarning('Weak Password', 'Password must contain at least one special character (!@#$%^&*)');
            return false;
        }

        // Confirm password
        if (password !== confirmPassword) {
            await showError('Password Mismatch', 'Passwords do not match');
            return false;
        }

        // Name validation
        if (!firstName || !lastName) {
            await showWarning('Missing Information', 'First name and last name are required');
            return false;
        }

        return true;
    }

    async function validateStep2() {
        const email = document.getElementById('email').value;
        const phone = document.getElementById('phone').value;
        const address = document.getElementById('address').value;
        const city = document.getElementById('city').value;
        const zipCode = document.getElementById('zipCode').value;

        // Email validation
        if (!/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(email)) {
            await showWarning('Invalid Email', 'Please enter a valid email address');
            return false;
        }

        // Required fields
        if (!phone || !address || !city || !zipCode) {
            await showWarning('Missing Information', 'All fields are required');
            return false;
        }

        return true;
    }

    async function validateStep3() {
        const frontId = document.getElementById('idCardFront').files[0];
        const backId = document.getElementById('idCardBack').files[0];
        const termsCheck = document.getElementById('termsCheck').checked;

        if (!frontId || !backId) {
            await showWarning('Missing Documents', 'Both front and back ID card images are required');
            return false;
        }

        if (!termsCheck) {
            await showWarning('Terms Agreement Required', 'You must agree to the Terms of Service and Privacy Policy');
            return false;
        }

        return true;
    }
});

function fileToBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => {
            const base64String = reader.result.split(',')[1];
            resolve(base64String);
        };
        reader.onerror = error => reject(error);
    });
}

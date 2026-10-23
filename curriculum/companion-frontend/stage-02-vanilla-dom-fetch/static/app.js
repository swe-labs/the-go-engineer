/**
 * The Go Engineer — Frontend Companion Track Stage 2
 * Modern Vanilla JavaScript, DOM Event Delegation & Asynchronous Fetch API
 */

(function () {
    "use strict";

    // DOM Elements Cache
    const searchInput = document.getElementById("search-input");
    const searchSpinner = document.getElementById("search-spinner");
    const booksContainer = document.getElementById("books-container");
    const booksCount = document.getElementById("books-count");
    const addBookForm = document.getElementById("add-book-form");
    const submitBtn = document.getElementById("submit-btn");
    const formAlert = document.getElementById("form-alert");

    // Local State
    let currentBooks = [];

    // =========================================================================
    // 1. Utilities: Debounce & Sanitization
    // =========================================================================

    /**
     * Prevents excessive API queries by delaying execution until typing pauses.
     * @param {Function} func 
     * @param {number} delayMs 
     */
    function debounce(func, delayMs) {
        let timeoutId;
        return function (...args) {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => func.apply(this, args), delayMs);
        };
    }

    /**
     * Escapes HTML special characters to prevent cross-site scripting (XSS).
     * @param {string} str 
     */
    function escapeHTML(str) {
        if (!str) return "";
        return String(str)
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    // =========================================================================
    // 2. Fetch API Network Calls
    // =========================================================================

    /**
     * Fetches books from the REST API matching an optional query string.
     * @param {string} query 
     */
    async function loadBooks(query = "") {
        searchSpinner.hidden = false;
        try {
            const url = query ? `/api/books?q=${encodeURIComponent(query)}` : "/api/books";
            const res = await fetch(url);
            if (!res.ok) {
                throw new Error(`HTTP error! status: ${res.status}`);
            }
            const data = await res.json();
            currentBooks = Array.isArray(data) ? data : [];
            renderBooks(currentBooks);
        } catch (err) {
            console.error("Failed to load books:", err);
            booksContainer.innerHTML = `
                <div class="empty-state error">
                    <p>Failed to connect to Book Catalog API: ${escapeHTML(err.message)}</p>
                </div>
            `;
        } finally {
            searchSpinner.hidden = true;
        }
    }

    /**
     * Deletes a book by ID.
     * @param {number} bookId 
     */
    async function deleteBook(bookId) {
        try {
            const res = await fetch(`/api/books/${bookId}`, {
                method: "DELETE"
            });

            if (res.status === 204) {
                // Optimistically remove from state and UI
                currentBooks = currentBooks.filter(b => b.id !== bookId);
                renderBooks(currentBooks);
                showAlert("success", "Book successfully removed from library.");
            } else if (res.status === 404) {
                showAlert("error", "Book not found or already deleted.");
                loadBooks(searchInput.value.trim());
            } else {
                throw new Error(`Server returned status ${res.status}`);
            }
        } catch (err) {
            console.error("Failed to delete book:", err);
            showAlert("error", `Delete failed: ${err.message}`);
        }
    }

    /**
     * Submits a new book to the REST API with field error handling.
     * @param {Object} payload 
     */
    async function createBook(payload) {
        clearFieldErrors();
        setFormSubmitting(true);

        try {
            const res = await fetch("/api/books", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(payload)
            });

            const data = await res.json();

            if (res.status === 201) {
                // Success: reset form and refresh list
                addBookForm.reset();
                showAlert("success", `Added "${escapeHTML(data.title)}" to catalog!`);
                await loadBooks(searchInput.value.trim());
            } else if (res.status === 400 && data.details) {
                // Validation Errors returned by Go backend
                renderValidationErrors(data.details);
                showAlert("error", data.error || "Please correct the errors in the form.");
            } else {
                showAlert("error", data.error || "An unexpected error occurred.");
            }
        } catch (err) {
            console.error("Create book network failure:", err);
            showAlert("error", `Network request failed: ${err.message}`);
        } finally {
            setFormSubmitting(false);
        }
    }

    // =========================================================================
    // 3. DOM Rendering & Manipulation
    // =========================================================================

    /**
     * Renders book cards in the catalog container.
     * @param {Array} books 
     */
    function renderBooks(books) {
        booksCount.textContent = books.length;

        if (books.length === 0) {
            booksContainer.innerHTML = `
                <div class="empty-state">
                    <p class="empty-title">No books found</p>
                    <p class="empty-subtitle">Try adjusting your search criteria or add a new title.</p>
                </div>
            `;
            return;
        }

        const cardsHTML = books.map(book => `
            <article class="book-card" data-id="${book.id}">
                <div class="book-card-main">
                    <div class="book-header">
                        <h3 class="book-title">${escapeHTML(book.title)}</h3>
                        <span class="book-price">$${Number(book.price).toFixed(2)}</span>
                    </div>
                    <p class="book-author">by ${escapeHTML(book.author)}</p>
                    <div class="book-meta">
                        <span class="meta-tag">ISBN: ${escapeHTML(book.isbn)}</span>
                        <span class="meta-tag">ID: #${book.id}</span>
                    </div>
                </div>
                <div class="book-card-actions">
                    <button 
                        type="button" 
                        class="btn-delete" 
                        data-action="delete" 
                        data-id="${book.id}"
                        data-title="${escapeHTML(book.title)}"
                        title="Delete book"
                    >
                        Delete
                    </button>
                </div>
            </article>
        `).join("");

        booksContainer.innerHTML = cardsHTML;
    }

    /**
     * Displays field-level validation errors under inputs.
     * @param {Array} errors 
     */
    function renderValidationErrors(errors) {
        errors.forEach(err => {
            const input = document.getElementById(`book-${err.field}`);
            const errorSpan = document.getElementById(`error-${err.field}`);
            if (input) {
                input.classList.add("input-error");
            }
            if (errorSpan) {
                errorSpan.textContent = err.message;
            }
        });
    }

    function clearFieldErrors() {
        document.querySelectorAll(".field-error").forEach(el => el.textContent = "");
        document.querySelectorAll(".input-error").forEach(el => el.classList.remove("input-error"));
    }

    function showAlert(type, message) {
        formAlert.className = `alert-box alert-${type}`;
        formAlert.textContent = message;
        formAlert.hidden = false;

        if (type === "success") {
            setTimeout(() => {
                formAlert.hidden = true;
            }, 4000);
        }
    }

    function setFormSubmitting(isSubmitting) {
        submitBtn.disabled = isSubmitting;
        const spinner = submitBtn.querySelector(".btn-spinner");
        const text = submitBtn.querySelector(".btn-text");
        if (spinner && text) {
            spinner.hidden = !isSubmitting;
            text.textContent = isSubmitting ? "Saving..." : "Add Book to Catalog";
        }
    }

    // =========================================================================
    // 4. Event Listeners & Delegation Patterns
    // =========================================================================

    // Debounced Search Input (300ms)
    searchInput.addEventListener("input", debounce(e => {
        loadBooks(e.target.value.trim());
    }, 300));

    // Event Delegation for Asynchronous Delete Buttons
    booksContainer.addEventListener("click", e => {
        const deleteBtn = e.target.closest('[data-action="delete"]');
        if (!deleteBtn) return;

        const id = parseInt(deleteBtn.dataset.id, 10);
        const title = deleteBtn.dataset.title;

        if (window.confirm(`Are you sure you want to remove "${title}"?`)) {
            deleteBook(id);
        }
    });

    // Add Book Form Submission
    addBookForm.addEventListener("submit", e => {
        e.preventDefault();

        const formData = new FormData(addBookForm);
        const payload = {
            title: formData.get("title").trim(),
            author: formData.get("author").trim(),
            isbn: formData.get("isbn").trim(),
            price: parseFloat(formData.get("price")) || 0
        };

        createBook(payload);
    });

    // Initial Load
    document.addEventListener("DOMContentLoaded", () => {
        loadBooks();
    });

})();

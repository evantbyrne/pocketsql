import * as Turbo from "@hotwired/turbo"

// After render event
// see: https://github.com/hotwired/turbo/issues/1289
// and: https://discuss.hotwired.dev/t/event-to-know-a-turbo-stream-has-been-rendered/1554/25
const afterRenderEvent = new Event("turbo:after-stream-render");
document.addEventListener("turbo:before-stream-render", (event) => {
  const originalRender = event.detail.render;
  event.detail.render = function (streamElement) {
    originalRender(streamElement);
    document.dispatchEvent(afterRenderEvent);
  }
});

const reload = function() {
	[...document.querySelectorAll("tr[data-modal-id]")].forEach(row => {
		const onSelect = () => {
			[...document.querySelectorAll('.modal:not(#' + row.getAttribute("data-modal-id") + ')')].forEach(modal => {
				modal.classList.remove("-active");
			})
			document.querySelector('.modal#' + row.getAttribute("data-modal-id")).classList.add("-active");
		};
		row.addEventListener("click", onSelect);
		row.addEventListener("keydown", event => {
			if (event.key === "Enter") {
				onSelect(event);
			}
		});
	});
	
	[...document.querySelectorAll('.modal_close,.shade')].forEach(node => {
		const onSelect = () => {
			[...document.querySelectorAll('.modal')].forEach(modal => {
				modal.classList.remove("-active");
			});
		}
		node.addEventListener("click", onSelect);
		node.addEventListener("keydown", event => {
			if (event.key === "Enter") {
				onSelect(event);
			}
		});
	});

	[...document.querySelectorAll('.query textarea')].forEach(node => {
		const resizeTextarea = () => {
			node.style.height = '';
			node.style.height = `${node.scrollHeight}px`;
		};
		resizeTextarea();
		node.addEventListener('input', resizeTextarea);
	});
};

document.addEventListener("turbo:load", reload);
document.addEventListener("turbo:after-stream-render", reload);

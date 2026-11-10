import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rails = el.querySelectorAll('path[d="M5 4v16"], path[d="M12 4v16"]');
  animate(rails, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const gate = el.querySelector('path[d^="M19 4"]');
  animate(gate, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

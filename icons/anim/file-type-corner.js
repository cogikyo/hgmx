import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M12 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const letter = el.querySelectorAll('path[d^="M3 16"], path[d="M7 14v8"]');
  animate(letter, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const foot = el.querySelector('path[d="M6 22h2"]');
  animate(foot, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}

import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pad = el.querySelectorAll('path[d^="M20 20"], path[d^="M9.656"]');
  animate(pad, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

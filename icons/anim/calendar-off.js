import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelectorAll('path[d^="M3.586"], path[d^="M8.656"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

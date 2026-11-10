import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ham = el.querySelectorAll('path[d^="M13.144"]');
  animate(ham, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const shine = el.querySelector('path[d^="m8.5"]');
  animate(shine, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}

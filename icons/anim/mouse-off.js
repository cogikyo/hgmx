import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M18.218"], path[d^="M19 13.343"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d="M22 22 2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

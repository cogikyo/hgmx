import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const wallet = el.querySelector('path[d^="M7 7"]');
  animate(wallet, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const clasp = el.querySelector('path[d="M17 14h.01"]');
  animate(clasp, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}

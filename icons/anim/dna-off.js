import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const strands = el.querySelectorAll('path[d^="M15 2"], path[d^="M2 15"], path[d^="M22 9"], path[d^="M9 22"]');
  animate(strands, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

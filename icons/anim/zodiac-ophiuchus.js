import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bearer = el.querySelector('path[d^="M6 3"]');
  animate(bearer, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const serpent = el.querySelector('path[d^="M3 10"]');
  animate(serpent, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

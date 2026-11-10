import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M20 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector('path[d="M2 10h20"]');
  animate(lid, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

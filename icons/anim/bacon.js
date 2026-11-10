import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rasher = el.querySelector('path[d^="M2.795"]');
  animate(rasher, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const fat = el.querySelector('path[d^="M4.5"]');
  animate(fat, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

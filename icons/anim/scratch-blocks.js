import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const blocks = el.querySelector('path[d^="M19 5"]');
  animate(blocks, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const joint = el.querySelector('path[d^="M21 12"]');
  animate(joint, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

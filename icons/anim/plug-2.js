import { animate, τmax } from "../motion.js";

export default function (el) {
  const prongs = el.querySelectorAll('path[d="M9 2v6"], path[d="M15 2v6"]');
  animate(prongs, { scaleY: [0, 1.25, 1] }, { duration: τmax });
}

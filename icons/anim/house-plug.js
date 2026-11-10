import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M8.5 21"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const prongs = el.querySelector("g");
  animate(prongs, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}

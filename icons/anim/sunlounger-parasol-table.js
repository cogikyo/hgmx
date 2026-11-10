import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const canopy = el.querySelector('path[d="M21 8H3l9-6Z"]');
  animate(canopy, { scaleX: [0.3, 1.15, 1] }, { duration: τmid });

  const table = el.querySelector('path[d="M8 13h8"]');
  animate(table, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}

import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const tiles = el.querySelectorAll("rect");
  animate(tiles, { translateX: [0, 4, -1, 0] }, { duration: τmax });

  const heads = el.querySelectorAll('path[d^="m7 14"], path[d^="m17 14"]');
  animate(heads, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}

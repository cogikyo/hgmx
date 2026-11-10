import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const tiles = el.querySelectorAll("rect");
  animate(tiles, { translateY: [0, 4, -1, 0] }, { duration: τmax });

  const heads = el.querySelectorAll('path[d^="m2 7"], path[d^="m2 17"]');
  animate(heads, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}

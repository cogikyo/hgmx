import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const tag = el.querySelector("g");
  animate(tag, { rotate: [0, -18, 14, -8, 0] }, { duration: τmax });

  const x = el.querySelectorAll('path[d^="m16.5 10.5"], path[d^="m21.5 10.5"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(τmin) });
}

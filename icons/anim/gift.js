import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelector('path[d^="M20 11"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector("g.lid");
  animate(lid, { translateY: [0, -4, 1, 0] }, { duration: τmid, delay: τmin });

  const bow = el.querySelector('path[d^="M7.5 7"]');
  animate(bow, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}

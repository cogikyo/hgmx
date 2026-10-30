import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector("rect");
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const x = [el.querySelector('path[d^="m14 13"]'), el.querySelector('path[d^="m10 13"]')];
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(τmin) });
}

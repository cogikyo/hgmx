import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelector("rect");
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stripes = ["M10 14", "m8 6", "m14 6"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(stripes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}

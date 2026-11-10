import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const loop = el.querySelector("g");
  animate(loop, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmax });

  const heads = ["m14 16", "M8.293", "m13.378"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(heads, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}

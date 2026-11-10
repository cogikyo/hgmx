import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const dashes = ["M9 3v2", "M9 9v1", "M9 14v1", "M9 19v2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(dashes, { translateX: [-6, 1, 0] }, { duration: τmax, delay: stagger(0.03) });
}

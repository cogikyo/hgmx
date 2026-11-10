import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const dashes = ["M15 3v2", "M15 9v1", "M15 14v1", "M15 19v2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(dashes, { translateX: [6, -1, 0] }, { duration: τmax, delay: stagger(0.03) });
}

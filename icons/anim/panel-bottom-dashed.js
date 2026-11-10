import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const dashes = ["M3 15h2", "M9 15h1", "M14 15h1", "M19 15h2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(dashes, { translateY: [6, -1, 0] }, { duration: τmax, delay: stagger(0.03) });
}

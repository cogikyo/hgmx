import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const top = ["M3 9h2", "M9 9h1", "M14 9h1", "M19 9h2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(top, { translateY: [-6, 1, 0] }, { duration: τmax, delay: stagger(0.03) });

  const bottom = ["M3 15h2", "M9 15h1", "M14 15h1", "M19 15h2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bottom, { translateY: [6, -1, 0] }, { duration: τmax, delay: stagger(0.03) });
}

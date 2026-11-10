import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const lens = el.querySelector("circle");
  animate(lens, { scale: [1, 0.85, 1] }, { duration: τmid });

  const rays = ["M5 7 3 5", "M9 6V3", "m13 7 2-2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(rays, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}

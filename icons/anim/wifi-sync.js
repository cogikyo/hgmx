import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const waves = ["M8.5", "M5", "M2 8"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const arrows = el.querySelector("g");
  animate(arrows, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmid, delay: τmin });
}

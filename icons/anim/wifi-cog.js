import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const waves = ["M8.5", "M5", "M2 8"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const gear = el.querySelector("g.gear");
  animate(gear, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}

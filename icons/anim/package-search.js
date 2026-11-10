import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const box = ["M12 22V12", "M21 1", "M3.29 7", "m7.5 4.27"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lens = el.querySelector("g");
  animate(lens, { scale: [1.35, 1] }, { duration: τmax });
  animate(lens, { translateX: [0, -3, 2, 0], translateY: [0, -3, -1, 0] }, { duration: τmax });
}

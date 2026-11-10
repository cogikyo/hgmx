import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cup = el.querySelector("g");
  animate(cup, { translateY: [-3, 0.5, 0] }, { duration: τmid });

  const steam = ["M6 4", "M12 4", "M18 4"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(steam, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}

import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const paths = el.querySelectorAll("path");

  const tailPaths = Array.from(paths).slice(1);
  animate(
    tailPaths,
    { scale: [0, 1.75, 1], opacity: [0, 1] },
    { duration: 0.3, delay: stagger(0.1) },
  );

  const mainLogo = paths[0];
  animate(mainLogo, { opacity: [0, 1] }, { duration: τmax });

  const group = el.querySelector("g");
  animate(group, { rotateZ: [0, -360], scale: [0.25, 1.08, 1] }, { duration: τmax });
}
